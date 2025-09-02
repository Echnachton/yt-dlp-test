export const getFile = async (id: string) => {
  const response = await fetch(`http://localhost:8080/api/v1/file/${id}`);
  return response.blob();
};
