"use client";
import React from "react";
import { SeedItem } from "../seedItem/SeedItem";
import { useRouter } from "next/navigation";
import { useState, useEffect, FC } from "react";
import { EyeInvisibleOutlined, EyeTwoTone } from "@ant-design/icons";
import { apiService } from "../axios/global.service";
import { Input, Space, Button, List, Card, Typography } from "antd";
import { LockOutlined } from "@ant-design/icons";
import useLoginHooks from "@/hooks/login";
import Cookies from "js-cookie";

export function LoginPage() {
  const { passwordVisible, setPasswordVisible } = useLoginHooks();
  const [inp, setInp] = useState("");
  //
  const [seed, setSeed] = useState<string[]>([
    "airport",
    "donate",
    "language",
    "disagree",
    "dumb",
    "access",
    "insect",
    "tribe",
    "ozone",
    "humor",
    "foot",
    "jealous",
    "much",
    "digital",
  ]);

  // const [seed, setSeed] = useState<string[]>(["melody","correct","brain","slide","flip","polar","asset","know","pencil","major","smile","vital","nominee","merge","addict"]);
  const [disabled, setDisabled] = useState(true);
  const router = useRouter();

  const onLogin = async () => {
    const result = await apiService.login(seed.join(" "));
    // if (result) {
    Cookies.set("accessToken", result.data.token);
    router.push("/organization");
    // }
  };
  const onRegister = async () => {
    const result = await apiService.register(inp);
    if (result) {
      Cookies.set("accessToken", result.data.token);
      router.push("/organization");
    }
  };
  // const getSeed = (event: React.ChangeEvent<HTMLInputElement>) => {
  //   setSeed(inp);
  // };
  console.log(inp);
  const { Text } = Typography;
  useEffect(() => {
    setDisabled(!(inp.length >= 4));
  }, [inp]);

  const onSubmitHandler = async () => {
    setSeed((prevState) => (prevState ? [...prevState, inp] : [inp]));
    setInp("");
  };

  return (
    <div className="flex relative  overflow-hidden flex-col w-2/3 h-3/4 items-center justify-center gap-10 bg-white border-solid border rounded-md border-neutral-300 text-neutral-500">
      <div className="w-full h-20    bg-[#1677FF] absolute top-0 flex items-center justify-center">
        <h1 className="text-white text-xl font-semibold">Log In</h1>
      </div>
      <div className="flex flex-col w-6/12  gap-3 items-start mt-20">
        <div>
          <Text type="secondary">
            Please enter by spaces 15 words always lowercase.
          </Text>
        </div>

        <Space.Compact style={{ width: "100%" }}>
          <Input.Password
            defaultValue="Combine input and button"
            value={inp}
            type="submit"
            size="large"
            suffix={<LockOutlined className="site-form-item-icon" />}
            addonBefore="Seed"
            visibilityToggle={{
              visible: !!passwordVisible,
              // onVisibleChange: setPasswordVisible,
            }}
            iconRender={(visible) =>
              visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />
            }
            placeholder="Enter your seed words in order"
            onInput={(event: any) => setInp(event.target.value)}
            //maxLength={8}
            minLength={4}
            //onChange={getSeed}
          />
          <Button
            size="large"
            type="primary"
            disabled={disabled}
            onClick={onSubmitHandler}
          >
            Submit
          </Button>
        </Space.Compact>
      </div>
      {
        <div
          className="flex   flex-row w-[700px] gap-3 content-box flex-wrap
      "
        >
          {seed &&
            seed.map((element: string, index: number) => (
              <SeedItem key={index} seed={element} />
            ))}
        </div>
      }
      <Button
        onClick={onLogin}
        style={{ width: "150px" }}
        type="primary"
        size="large"
      >
        Login
      </Button>
      <Button
        onClick={onRegister}
        style={{ width: "150px" }}
        type="primary"
        size="large"
      >
        Join
      </Button>
    </div>
  );
}
