import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { type DownloadResponse, getDownloads } from "./queries/download/get";
import { postDownload } from "./queries/download/post";

const headers = ["ID", "URL", "Owner ID"];

export const Page = () => {
  const [url, setUrl] = useState("");

  const { data, isLoading, error: queryError } = useQuery<DownloadResponse>({
    queryKey: ["downloads"],
    queryFn: getDownloads,
  });

  const { mutate, isPending, error: mutationError } = useMutation({
    mutationFn: postDownload,
    onSuccess: (data) => {
      console.log("Download job queued successfully:", data);
      setUrl("");
    },
    onError: (error) => {
      console.error("Failed to queue download job:", error);
    },
  });

  const error = queryError || mutationError;

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!url) {
      alert("Please enter a URL");
      return;
    }
    mutate(url);
  };

  console.log(data);

  return (
    <div>
      <form onSubmit={handleSubmit}>
        <label htmlFor="url">URL</label>
        <input
          type="text"
          id="url"
          style={{ border: "solid 1px black" }}
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
        <button type="submit" disabled={isPending}>Download</button>
      </form>
      {isLoading && <p>Loading...</p>}

      {error && <p>Error: {error.message}</p>}

      {data && (
        <Table>
          <TableHeader>
            <TableRow>
              {headers.map((header) => (
                <TableHead key={header}>{header}</TableHead>
              ))}
            </TableRow>
          </TableHeader>

          <TableBody>
            {data?.videos?.map((download) => (
              <TableRow key={download.id}>
                <TableCell>{download.id}</TableCell>
                <TableCell>{download.url}</TableCell>
                <TableCell>{download.ownerId}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </div>
  );
};
