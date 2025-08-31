import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Spinnner } from "@/components/ui/spinnner";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { type DownloadResponse, getDownloads } from "@/queries/download/get";
import { postDownload } from "@/queries/download/post";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Bug } from "lucide-react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

const headers = ["ID", "URL", "Owner ID"];

const formSchema = z.object({
  url: z
    .url("Please enter a valid URL")
    .min(1, "URL is required")
    .refine(
      (url) => {
        const youtubeRegex =
          /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/.+/;
        return youtubeRegex.test(url);
      },
      {
        message: "Please enter a valid YouTube URL",
      },
    ),
});

type FormData = z.infer<typeof formSchema>;

export const Page = () => {
  const queryClient = useQueryClient();
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<FormData>({
    resolver: zodResolver(formSchema),
  });

  const { data, isLoading, error: queryError } = useQuery<DownloadResponse>({
    queryKey: ["downloads"],
    queryFn: getDownloads,
  });

  const { mutate, isPending, error: mutationError } = useMutation({
    mutationFn: postDownload,
    onSuccess: () => {
      reset();
      queryClient.invalidateQueries({ queryKey: ["downloads"] });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  const error = queryError || mutationError;

  const onSubmit = (data: FormData) => {
    mutate(data.url);
  };

  return (
    <div className="flex gap-1 flex-col p-4">
      <Card>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="flex gap-1 flex-col">
              <div className="flex gap-1">
                <Input
                  type="text"
                  id="url"
                  style={{ border: "solid 1px black" }}
                  {...register("url")}
                  placeholder="Enter a YouTube URL"
                />

                <Button type="submit" disabled={isPending}>
                  {isPending && <Spinnner />}
                  Download
                </Button>
              </div>

              {errors.url && (
                <span className="text-red-500 text-sm">
                  {errors.url.message}
                </span>
              )}
            </div>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                {headers.map((header) => (
                  <TableHead key={header}>{header}</TableHead>
                ))}
              </TableRow>
            </TableHeader>

            <TableBody>
              {isLoading && (
                <TableRow>
                  <TableCell colSpan={headers.length}>
                    <div className="flex items-center justify-center">
                      <Spinnner />
                    </div>
                  </TableCell>
                </TableRow>
              )}

              {error && (
                <TableRow>
                  <TableCell colSpan={headers.length}>
                    <Alert variant="destructive">
                      <Bug />

                      <AlertTitle>Error</AlertTitle>

                      <AlertDescription>
                        {error.message}
                      </AlertDescription>
                    </Alert>
                  </TableCell>
                </TableRow>
              )}

              {data?.videos?.map((download) => (
                <TableRow key={download.id}>
                  <TableCell>{download.id}</TableCell>
                  <TableCell>{download.url}</TableCell>
                  <TableCell>{download.owner_id}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};
