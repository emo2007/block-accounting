"use client";
import { Button } from "antd";
import { useRouter } from "next/navigation";
import { ArrowLeftOutlined } from "@ant-design/icons";
export function BackBtn() {
  const router = useRouter();
  const goBack = () => {
    router.back();
  };
  return (
    <Button
      className="flex items-center"
      type="text"
      size="large"
      onClick={goBack}
    >
      <ArrowLeftOutlined />
    </Button>
  );
}
