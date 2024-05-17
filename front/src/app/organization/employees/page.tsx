"use client";
import React from "react";
import { EmployeePage } from "./Employee";

export default function Home() {
  return (
    <div className="flex w-full h-screen items-center justify-center bg-white">
      <EmployeePage />
    </div>
  );
}
