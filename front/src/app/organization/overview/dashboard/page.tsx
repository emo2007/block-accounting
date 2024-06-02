"use client";
import React from "react";
import { OrgProfile } from "./[id]";
import { useSearchParams } from "next/navigation";

export default function Home() {
  return (
    <div className="flex flex-row w-full h-screen bg-slate-50 gap-5">
      <OrgProfile />
    </div>
  );
}
