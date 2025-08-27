export const postDownload = async (url: string) => {
  const response = await fetch("http://localhost:8080/api/v1/download", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "test_user": "ui_user", // Add the required header
    },
    body: JSON.stringify({ url }),
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(
      errorData.message || `HTTP error! status: ${response.status}`,
    );
  }

  return response.json();
};
