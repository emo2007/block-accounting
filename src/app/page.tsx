"use client";
import { LoginPage } from "./login/LoginPage";

export default function Home() {
  return (
    <div className="flex w-full h-screen items-center justify-center bg-white">
      <LoginPage />
    </div>
  );
}
