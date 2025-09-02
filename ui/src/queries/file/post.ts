type PostFileRequest = {
  id: string;
};

export const postFile = async (request: PostFileRequest) => {
  const response = await fetch(`http://localhost:8080/api/v1/file`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  });

  return response.blob();
};
