import { api } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { useAppContext } from "../context/AppContext";

export const getListPromptTypes = async (token: string, page: number) => {
	const response = await api.listPromptTypes(token, page);
	console.log("Fetched Prompt Types:", response); // Log the response for debugging
	return response;
};

export const useFetchPromptType = (page: number) => {
	const { token } = useAppContext();
	return useQuery({
		queryFn: () => getListPromptTypes(token!, page),
		queryKey: ["promptTypes", page],
		enabled: !!token,
	});
};
