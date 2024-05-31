import { useState } from "react";

export default function useLoginHooks() {
  const [passwordVisible, setPasswordVisible] = useState([]);

  // const getSeed = (string: string) => {
  //   // if (seed.length < 12) {
  //   setSeed((prev) => [...prev, string]);
  // };

  return { passwordVisible, setPasswordVisible };
}
