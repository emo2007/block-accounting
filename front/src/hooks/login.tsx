import { useState } from "react";

export default function useLoginHooks() {
  const [passwordVisible, setPasswordVisible] = useState([]);
  const [seed, setSeed] = useState<string[]>([]);

  const getSeed = (string: string) => {
    if (seed.length < 12) {
      setSeed((prev) => [...prev, string]);
    }
  };
  console.log(seed);

  return { passwordVisible, setPasswordVisible, seed, getSeed };
}
