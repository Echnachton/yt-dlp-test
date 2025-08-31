import { LoaderCircle } from "lucide-react";
import { motion } from "motion/react";

export const Spinnner = () => {
  return (
    <motion.div
      animate={{ rotate: 360 }}
      transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
      className="inline-flex items-center justify-center"
    >
      <LoaderCircle />
    </motion.div>
  );
};
