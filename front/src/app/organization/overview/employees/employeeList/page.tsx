"use client";
import React from "react";
import { EmployeeList } from "./EmployeeList";

export default function Home() {
  return (
    <div className="flex w-full h-screen items-center justify-center bg-slate-50">
      <EmployeeList />
    </div>
  );
}
