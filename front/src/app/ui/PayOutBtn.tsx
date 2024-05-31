"use client";
import { Button } from "antd";
import { useRouter } from "next/navigation";
export function PayOutBtn() {
  const router = useRouter();
  const onPayHandler = () => {
    router.push("/organization/overview/pending");
  };
  return (
    <Button onClick={onPayHandler} type="primary">
      Pay Out
    </Button>
  );
}
