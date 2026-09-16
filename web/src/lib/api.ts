import axios, { AxiosError, AxiosResponse } from "axios";

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type User = {
	id: string;
	email: string;
	plan: string;
	credit_balance: number;
	created_at: string;
	updated_at: string;
};

export type Job = {
	id: string;
	user_id: string;
	status: string;
	prompt_input: Record<string, string>;
	assembled_prompt: string;
	image_url?: string;
	model_name: string;
	duration_seconds: number;
	credits_charged: number;
	gcs_video_url?: string | null;
	error_message?: string | null;
	created_at: string;
	updated_at: string;
};

export type CreditTransaction = {
	id: string;
	amount: number;
	action: string;
	job_id?: string;
	created_at: string;
};

export type PromptType = {
	id: string;
	name: string;
	type: string;
	description: string;
	icon: string;
	background_class: string;
	is_deleted?: boolean;
	created_at: string;
};

export type APIResponse<T> = {
	success: boolean;
	data?: T;
	error?: string;
	message?: string;
};

const apiClient = axios.create({
	baseURL: API_BASE,
});

apiClient.interceptors.response.use(
	undefined,
	(error: AxiosError<APIResponse<unknown>>) => {
		const errMessage = error.response?.data?.error || error.message;
		
		// Only redirect to login if the error is a 401 Unauthorized
		if (error.response?.status === 401) {
			location.href = "/login";
		}
		
		return Promise.reject(new Error(errMessage));
	},
);

async function unwrap<T>(request: Promise<AxiosResponse<APIResponse<T>>>): Promise<T> {
	const response = await request;
	const body = response.data;

	if (!body.success || body.data === undefined) {
		throw new Error(body.error || "Invalid API response");
	}

	return body.data;
}

export const api = {
	// Auth — response is already unwrapped by interceptor to { user, token }
	register: (email: string, password: string) =>
		unwrap(apiClient.post<APIResponse<{ user: User; token: string }>>("/api/v1/auth/register", {
			email,
			password,
		})),

	login: (email: string, password: string) =>
		unwrap(apiClient.post<APIResponse<{ user: User; token: string }>>("/api/v1/auth/login", {
			email,
			password,
		})),

	// User — unwrapped directly to User
	getMe: (token: string) =>
		unwrap(apiClient.get<APIResponse<User>>("/api/v1/me", {
			headers: { Authorization: `Bearer ${token}` },
		})),

	getCreditHistory: (token: string) =>
		unwrap(apiClient.get<APIResponse<CreditTransaction[]>>("/api/v1/me/credits", {
			headers: { Authorization: `Bearer ${token}` },
		})),

	// Jobs
	createJob: (
		token: string,
		promptInput: Record<string, string>,
		imageFile?: File,
	) => {
		const formData = new FormData();
		formData.append("prompt_input", JSON.stringify(promptInput));
		if (imageFile) {
			formData.append("image", imageFile);
		}
		return unwrap(apiClient.post<APIResponse<{ job_id: string; status: string }>>(
			"/api/v1/jobs",
			formData,
			{
				headers: {
					Authorization: `Bearer ${token}`,
				},
			},
		));
	},

	getJob: (token: string, jobId: string) =>
		unwrap(apiClient.get<APIResponse<Job>>(`/api/v1/jobs/${jobId}`, {
			headers: { Authorization: `Bearer ${token}` },
		})),

	listJobs: (token: string, page = 1) =>
		unwrap(apiClient.get<APIResponse<Job[]>>(`/api/v1/jobs?page=${page}&limit=20`, {
			headers: { Authorization: `Bearer ${token}` },
		})),

	listPromptTypes: (token: string, page = 1) =>
		unwrap(apiClient.get<APIResponse<{ items: PromptType[]; pagination: unknown }>>(`/api/v1/prompt-types?page=${page}&limit=100`, {
			headers: { Authorization: `Bearer ${token}` },
		})),
};
