"use client";
import React from "react";
import { OrgCreatePage } from "../orgCreate/OrgCreatePage";

export default function organization() {
  return (
    <div className="flex w-full h-screen items-center justify-center bg-white">
      <OrgCreatePage />
    </div>
  );
}
