import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "sonner";
import { Page } from "./Page";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Toaster
        position="bottom-center"
        richColors
      />

      <Page />
    </QueryClientProvider>
  );
}

export default App;
