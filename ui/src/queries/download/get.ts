export type Download = {
  id: number;
  url: string;
  ownerId: string;
};

export type DownloadResponse = {
  videos: Download[];
};

export const getDownloads = async (): Promise<DownloadResponse> => {
  const response = await fetch("http://localhost:8080/api/v1/download");
  return response.json();
};
