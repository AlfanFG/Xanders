"use client";

import { useState, useRef, ChangeEvent } from "react";
import { useRouter } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "../lib/api";
import { useAppContext } from "./context/AppContext";
import { useFetchPromptType } from "./hooks/useFetchPromptType";
import { PromptType } from "../lib/api";

type GenerationStatus =
	| "idle"
	| "pending"
	| "in_queue"
	| "processing"
	| "completed"
	| "failed";

const isGenerationStatus = (value: string): value is GenerationStatus =>
	["pending", "in_queue", "processing", "completed", "failed"].includes(value);

const errorMessageFrom = (error: unknown) =>
	error instanceof Error ? error.message : "Failed to create job.";

export default function Home() {
	const router = useRouter();
	const { token, user } = useAppContext();
	const queryClient = useQueryClient();
	const fileInputRef = useRef<HTMLInputElement>(null);
	const page = 1;
	const [step, setStep] = useState(1);
	const [imageFile, setImageFile] = useState<File | null>(null);
	const [imagePreview, setImagePreview] = useState<string | null>(null);

	const [form, setForm] = useState({
		template: "",
		lighting: "",
		camera: "",
		audio: "",
		customPrompt: "",
		personModel: "",
	});

	const [currentJobId, setCurrentJobId] = useState<string | null>(null);
	const [status, setStatus] = useState<GenerationStatus>("idle");
	const [videoUrl, setVideoUrl] = useState("");
	const [errorMessage, setErrorMessage] = useState("");
	const creditRequestHref = `mailto:alfanfaturahman10@gmail.com?subject=${encodeURIComponent("Xanders credit request")}&body=${encodeURIComponent(`Hello, I would like to request more Xanders video-generation credits.\n\nAccount: ${user?.email ?? ""}\nCurrent balance: ${user?.credit_balance ?? 0}`)}`;
	const { data } = useFetchPromptType(page);
	const items = data?.items || [];
	const persons = items.filter((t: PromptType) => t.type === "person_model");
	const templates = items.filter((t: PromptType) => t.type === "template");
	const lightings = items.filter((t: PromptType) => t.type === "lighting");
	const cameras = items.filter((t: PromptType) => t.type === "camera");
	const audios = items.filter((t: PromptType) => t.type === "audio");
	const selectedForm = {
		...form,
		personModel: form.personModel || persons[0]?.id || "",
		template: form.template || templates[0]?.id || "",
		lighting: form.lighting || lightings[0]?.id || "",
		camera: form.camera || cameras[0]?.id || "",
		audio: form.audio || audios[0]?.id || "",
	};

	const updateForm = (key: keyof typeof form, value: string) => {
		setForm((prev) => ({ ...prev, [key]: value }));
	};

	const handleNext = () => setStep((s) => Math.min(s + 1, 6));
	const handlePrev = () => setStep((s) => Math.max(s - 1, 1));

	const handleImageUpload = (e: ChangeEvent<HTMLInputElement>) => {
		if (e.target.files && e.target.files[0]) {
			const file = e.target.files[0];
			setImageFile(file);
			setImagePreview(URL.createObjectURL(file));
		}
	};

	const createJobMutation = useMutation({
		mutationFn: () => api.createJob(token!, selectedForm, imageFile || undefined),
		onSuccess: (data) => {
			setCurrentJobId(data.job_id);
			setStatus("in_queue");
			queryClient.invalidateQueries({ queryKey: ["user", token] });
		},
		onError: (err: unknown) => {
			setStatus("failed");
			setErrorMessage(errorMessageFrom(err));
		},
	});

	useQuery({
		queryKey: ["job", currentJobId],
		queryFn: async () => {
			const data = await api.getJob(token!, currentJobId!);
			const jobStatus = isGenerationStatus(data.status) ? data.status : "failed";
			setStatus(jobStatus);

			if (jobStatus === "completed") {
				setVideoUrl(data.gcs_video_url ?? "");
				setCurrentJobId(null);
			} else if (jobStatus === "failed") {
				setErrorMessage(data?.error_message || "Video generation failed");
				setCurrentJobId(null);
			}
			return data;
		},
		enabled: !!currentJobId && !!token,
		refetchInterval: (query) => {
			const st = query.state.data?.status;
			return st === "completed" || st === "failed" ? false : 10000;
		},
	});

	const handleGenerate = () => {
		if (!token) {
			router.push("/login");
			return;
		}
		setStatus("pending");
		setVideoUrl("");
		setErrorMessage("");
		createJobMutation.mutate();
	};

	return (
		<div className="min-h-screen bg-black text-white p-6 pb-20 md:p-12 font-sans selection:bg-primary/30">
			<header className="max-w-6xl mx-auto mb-12 text-center md:text-left pt-16">
				<h1 className="text-4xl md:text-6xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-cyan-400 via-blue-500 to-purple-600 mb-4 tracking-tight">
					Create Cinematic Magic
				</h1>
				<p className="text-gray-400 text-lg max-w-2xl">
					Upload a reference image or pick a style, and let the AI direct your
					next masterpiece.
				</p>
			</header>

			<div className="max-w-6xl mx-auto grid grid-cols-1 lg:grid-cols-12 gap-8 relative">
				{/* BUILDER SECTION */}
				<div className="lg:col-span-7">
					<div className="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-6 md:p-8 shadow-2xl relative overflow-hidden">
						{/* Progress Stepper */}
						<div className="flex justify-between items-center mb-10 relative">
							<div className="absolute left-0 top-1/2 -translate-y-1/2 w-full h-1 bg-white/5 rounded-full -z-10"></div>
							{[
								{ num: 1, label: "Visuals" },
								{ num: 2, label: "Model" },
								{ num: 3, label: "Template" },
								{ num: 4, label: "Style" },
								{ num: 5, label: "Audio" },
								{ num: 6, label: "Review" },
							].map((s) => (
								<div key={s.num} className="flex flex-col items-center gap-2">
									<div
										className={`w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm transition-all duration-300 ${
											step >= s.num
												? "bg-gradient-to-br from-blue-500 to-purple-600 text-white shadow-[0_0_15px_rgba(139,92,246,0.5)]"
												: "bg-white/10 text-gray-500"
										}`}
									>
										{s.num}
									</div>
									<span
										className={`text-xs font-medium hidden md:block ${step >= s.num ? "text-white" : "text-gray-500"}`}
									>
										{s.label}
									</span>
								</div>
							))}
						</div>

						{/* STEP 1: VISUAL REFERENCE (IMAGE UPLOAD) */}
						{step === 1 && (
							<div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
								<h2 className="text-2xl font-bold">Visual Reference</h2>
								<p className="text-gray-400 text-sm">
									Upload an image to guide the AI&apos;s generation (Optional but
									recommended for best results).
								</p>

								<div
									className={`border-2 border-dashed rounded-2xl p-10 text-center cursor-pointer transition-all ${imagePreview ? "border-purple-500 bg-purple-500/5" : "border-white/20 hover:border-white/50 hover:bg-white/5"}`}
									onClick={() => fileInputRef.current?.click()}
								>
									{imagePreview ? (
										<div className="relative w-full aspect-video rounded-xl overflow-hidden">
											<img
												src={imagePreview}
												alt="Preview"
												className="object-cover w-full h-full"
											/>
											<div className="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity">
												<span className="text-white font-medium">
													Click to change image
												</span>
											</div>
										</div>
									) : (
										<div className="flex flex-col items-center justify-center gap-4 py-8">
											<span className="text-5xl">📸</span>
											<div>
												<p className="text-lg font-medium">
													Click or drag image here
												</p>
												<p className="text-sm text-gray-500">
													JPG, PNG up to 5MB
												</p>
											</div>
										</div>
									)}
									<input
										type="file"
										ref={fileInputRef}
										onChange={handleImageUpload}
										accept="image/png, image/jpeg"
										className="hidden"
									/>
								</div>

								<div>
									<label className="block text-sm font-medium text-gray-400 mb-2">
										Additional Creative Direction (Optional)
									</label>
									<textarea
										value={selectedForm.customPrompt}
										onChange={(e) => updateForm("customPrompt", e.target.value)}
										placeholder="Add details the video must include, such as a new action, setting, or camera moment..."
										className="w-full bg-black/40 border border-white/10 rounded-xl p-4 text-white placeholder-gray-600 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500 transition-all h-24 resize-none"
									/>
								</div>
							</div>
						)}

						{/* STEP 2: PERSON MODEL */}
						{step === 2 && (
							<div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
								<div>
									<h2 className="text-2xl font-bold">Person Model</h2>
									<p className="text-gray-400 text-sm mt-1">
										Choose the type of person to feature in your advertising
										video.
									</p>
								</div>
								<div className="grid grid-cols-1 md:grid-cols-2 gap-3">
									{persons.map((p: PromptType) => (
										<div
											key={p.id}
											onClick={() => updateForm("personModel", p.id)}
											className={`p-4 rounded-2xl cursor-pointer border transition-all duration-200 flex items-center gap-4 group ${
											selectedForm.personModel === p.id
													? "bg-gradient-to-r from-blue-900/40 to-purple-900/40 border-purple-500 shadow-[0_0_20px_rgba(168,85,247,0.15)]"
													: "bg-black/40 border-white/10 hover:border-white/30 hover:bg-white/5"
											}`}
										>
											<span
												className={`text-2xl w-12 h-12 rounded-xl flex items-center justify-center flex-shrink-0 bg-gradient-to-br ${p.background_class || "from-slate-700 to-slate-900"} group-hover:scale-110 transition-transform duration-300`}
											>
												{p.icon}
											</span>
											<div className="flex-1 min-w-0">
												<h3 className="font-bold text-base leading-tight">
													{p.name}
												</h3>
												<p className="text-xs text-gray-400 mt-0.5 leading-snug">
													{p.description}
												</p>
											</div>
											{selectedForm.personModel === p.id && (
												<span className="text-purple-400 text-lg flex-shrink-0">
													✓
												</span>
											)}
										</div>
									))}
								</div>
							</div>
						)}

						{/* STEP 3: TEMPLATE */}
						{step === 3 && (
							<div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
								<div>
									<h2 className="text-2xl font-bold">Industry Template</h2>
									<p className="text-gray-400 text-sm mt-1">
										Select the product category for your advertisement.
									</p>
								</div>
								<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
									{templates.map((t: PromptType) => (
										<div
											key={t.id}
											onClick={() => updateForm("template", t.id)}
											className={`p-5 rounded-2xl cursor-pointer border transition-all duration-200 flex items-center gap-4 ${
											selectedForm.template === t.id
													? "bg-gradient-to-r from-blue-900/40 to-purple-900/40 border-purple-500 shadow-[0_0_20px_rgba(168,85,247,0.15)]"
													: "bg-black/40 border-white/10 hover:border-white/30"
											}`}
										>
											<span className="text-3xl bg-white/5 w-12 h-12 rounded-xl flex items-center justify-center">
												{t.icon}
											</span>
											<div>
												<h3 className="font-bold text-lg">{t.name}</h3>
												<p className="text-sm text-gray-400">{t.description}</p>
											</div>
										</div>
									))}
								</div>
							</div>
						)}

						{/* STEP 4: STYLE */}
						{step === 4 && (
							<div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
								<div>
									<h2 className="text-2xl font-bold mb-4">Lighting & Mood</h2>
									<div className="grid grid-cols-2 md:grid-cols-4 gap-4">
										{lightings.map((l: PromptType) => (
											<div
												key={l.id}
												onClick={() => updateForm("lighting", l.id)}
												className={`rounded-2xl cursor-pointer border overflow-hidden transition-all duration-200 flex flex-col group ${
													selectedForm.lighting === l.id
														? "border-purple-500 shadow-[0_0_15px_rgba(168,85,247,0.2)]"
														: "border-white/10 hover:border-white/30"
												}`}
											>
												<div
													className={`h-20 bg-gradient-to-br ${l.background_class || "from-slate-800 to-black"} flex items-center justify-center text-3xl group-hover:scale-105 transition-transform duration-500`}
												>
													{l.icon}
												</div>
												<div className="p-3 text-center bg-black/60 backdrop-blur-md font-medium text-sm">
													{l.name}
												</div>
											</div>
										))}
									</div>
								</div>
							</div>
						)}

						{/* STEP 5: AUDIO */}
						{step === 5 && (
							<div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
								<h2 className="text-2xl font-bold">Background Audio</h2>
								<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
									{audios.map((a: PromptType) => (
										<div
											key={a.id}
											onClick={() => updateForm("audio", a.id)}
											className={`p-5 rounded-2xl cursor-pointer border transition-all duration-200 flex items-center gap-4 ${
											selectedForm.audio === a.id
													? "bg-gradient-to-r from-blue-900/40 to-purple-900/40 border-purple-500 shadow-[0_0_20px_rgba(168,85,247,0.15)]"
													: "bg-black/40 border-white/10 hover:border-white/30"
											}`}
										>
											<span className="text-3xl bg-blue-500/20 text-blue-400 w-12 h-12 rounded-xl flex items-center justify-center">
												{a.icon}
											</span>
											<div>
												<h3 className="font-bold text-lg">{a.name}</h3>
											</div>
										</div>
									))}
								</div>
							</div>
						)}

						{/* STEP 6: REVIEW */}
						{step === 6 && (
							<div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
								<h2 className="text-2xl font-bold">Review & Generate</h2>
								<div className="bg-black/40 border border-white/10 rounded-2xl p-6 space-y-4">
									{imagePreview && (
										<div className="flex justify-between items-center border-b border-white/5 pb-4">
											<span className="text-gray-400">Reference Image</span>
											<img
												src={imagePreview}
												className="h-12 rounded bg-white/10"
												alt="Ref"
											/>
										</div>
									)}
									<div className="flex justify-between items-center border-b border-white/5 pb-4">
										<span className="text-gray-400">Person Model</span>
										<span className="font-bold flex items-center gap-2">
											{
												persons.find(
													(p: PromptType) => p.id === selectedForm.personModel,
												)?.icon
											}{" "}
											{persons.find(
												(p: PromptType) => p.id === selectedForm.personModel,
											)?.name || "N/A"}
										</span>
									</div>
									<div className="flex justify-between items-center border-b border-white/5 pb-4">
										<span className="text-gray-400">Template</span>
										<span className="font-bold">
											{templates.find((t: PromptType) => t.id === selectedForm.template)
												?.name || "N/A"}
										</span>
									</div>
									<div className="flex justify-between items-center border-b border-white/5 pb-4">
										<span className="text-gray-400">Lighting</span>
										<span className="font-bold">
											{lightings.find((l: PromptType) => l.id === selectedForm.lighting)
												?.name || "N/A"}
										</span>
									</div>
									<div className="flex justify-between items-center pb-2">
										<span className="text-gray-400">Cost</span>
										<span className="font-bold text-green-400 bg-green-400/10 px-3 py-1 rounded-full">
											5 Credits
										</span>
									</div>
								</div>
							</div>
						)}

						{/* NAVIGATION */}
						<div className="mt-10 flex justify-between items-center pt-6 border-t border-white/10">
							{step > 1 ? (
								<button
									className="px-6 py-3 rounded-xl font-medium text-white bg-white/5 hover:bg-white/10 transition-colors"
									onClick={handlePrev}
									disabled={
										status !== "idle" &&
										status !== "completed" &&
										status !== "failed"
									}
								>
									Back
								</button>
							) : (
								<div></div>
							)}

							{step < 6 ? (
								<button
									className="px-8 py-3 rounded-xl font-bold text-white bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 shadow-lg hover:shadow-purple-500/25 transition-all"
									onClick={handleNext}
								>
									Next Step
								</button>
							) : (
								<button
									className={`px-8 py-3 rounded-xl font-bold text-white shadow-lg transition-all ${
										status === "pending" ||
										status === "in_queue" ||
										status === "processing"
											? "bg-gray-700 cursor-not-allowed opacity-70"
											: "bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-500 hover:to-purple-500 hover:shadow-purple-500/25"
									}`}
									onClick={handleGenerate}
									disabled={
										status === "pending" ||
										status === "in_queue" ||
										status === "processing"
									}
								>
									{status === "idle" ||
									status === "completed" ||
									status === "failed"
										? "🎬 Generate Video"
										: "Processing..."}
								</button>
							)}
						</div>
					</div>
				</div>

				{/* PREVIEW SECTION (RIGHT SIDE) */}
				<div className="lg:col-span-5 h-[600px] lg:h-auto sticky top-24">
					<div className="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl h-full shadow-2xl overflow-hidden flex flex-col relative group">
						{/* IDLE */}
						{status === "idle" && (
							<div className="flex-1 flex flex-col items-center justify-center p-8 text-center bg-gradient-to-br from-black/50 to-transparent">
								<span className="text-6xl mb-6 opacity-50 group-hover:scale-110 transition-transform duration-500">
									✨
								</span>
								<p className="text-xl font-medium text-gray-300">
									Your masterpiece awaits
								</p>
								<p className="text-sm text-gray-500 mt-2">
									Configure settings to the left and click generate.
								</p>
							</div>
						)}

						{/* PROCESSING */}
						{(status === "pending" ||
							status === "in_queue" ||
							status === "processing") && (
							<div className="flex-1 flex flex-col items-center justify-center p-8 text-center bg-gradient-to-br from-blue-900/20 to-purple-900/20">
								<div className="w-16 h-16 border-4 border-purple-500/30 border-t-purple-500 rounded-full animate-spin mb-8"></div>
								<h3 className="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-400 mb-2">
									{status === "pending" && "Initializing request..."}
									{status === "in_queue" && "In Queue for AI Engine"}
									{status === "processing" && "Rendering AI Video..."}
								</h3>
								<p className="text-gray-400">
									Estimated time: 1-2 minutes. You can safely leave this page.
								</p>
							</div>
						)}

						{/* COMPLETED */}
						{status === "completed" && (
							<div className="w-full h-full flex flex-col bg-black">
								<video
									src={videoUrl}
									controls
									autoPlay
									loop
									className="w-full h-full object-cover"
								/>
								<div className="absolute bottom-0 left-0 w-full p-6 bg-gradient-to-t from-black via-black/80 to-transparent flex gap-4 opacity-0 group-hover:opacity-100 transition-opacity">
									<a
										href={videoUrl}
										target="_blank"
										rel="noopener noreferrer"
										className="flex-1"
									>
										<button className="w-full py-3 rounded-xl bg-white text-black font-bold hover:bg-gray-200 transition-colors">
											Download 4K
										</button>
									</a>
									<button
										className="px-6 py-3 rounded-xl bg-white/20 text-white font-medium hover:bg-white/30 backdrop-blur transition-colors"
										onClick={() => setStatus("idle")}
									>
										Reset
									</button>
								</div>
							</div>
						)}

						{/* FAILED */}
						{status === "failed" && (
							<div className="flex-1 flex flex-col items-center justify-center p-8 text-center">
								<span className="text-5xl mb-6 grayscale">⚠️</span>
								<p className="text-red-400 font-medium mb-6 bg-red-400/10 p-4 rounded-xl">
									{errorMessage}
								</p>
								{errorMessage.toLowerCase().includes("insufficient credits") && (
									<a
										href={creditRequestHref}
										className="mb-4 px-6 py-2 rounded-xl bg-gradient-to-r from-blue-600 to-purple-600 font-bold text-white hover:from-blue-500 hover:to-purple-500 transition-colors"
									>
										Request More Credits
									</a>
								)}
								<button
									className="px-6 py-2 rounded-xl bg-white/10 hover:bg-white/20 transition-colors"
									onClick={() => setStatus("idle")}
								>
									Try Again
								</button>
							</div>
						)}
					</div>
				</div>
			</div>
		</div>
	);
}
