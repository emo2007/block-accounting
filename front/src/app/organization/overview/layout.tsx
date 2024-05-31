"use client";
import React, { useState } from "react";
import Link from "next/link";
import { BackBtn } from "@/app/ui/BackBtn";
import {
  PieChartOutlined,
  KeyOutlined,
  UsergroupAddOutlined,
  FileProtectOutlined,
  WalletOutlined,
} from "@ant-design/icons";
import type { MenuProps } from "antd";
import { Menu } from "antd";

export default function Layout({ children }: { children: React.ReactNode }) {
  type MenuItem = Required<MenuProps>["items"][number];
  const [collapsed, setCollapsed] = useState(true);
  const items: MenuItem[] = [
    {
      key: "1",
      icon: <PieChartOutlined />,
      label: (
        <Link href="http://localhost:3000/organization/overview/dashboard">
          Overview
        </Link>
      ),
    },
    {
      key: "2",
      icon: <FileProtectOutlined />,

      children: [
        {
          key: "5",
          label: (
            <Link href="http://localhost:3000/organization/overview/license">
              Licenses
            </Link>
          ),
        },
        {
          key: "6",
          label: (
            <Link href="http://localhost:3000/organization/overview/agreement">
              Agreement
            </Link>
          ),
        },
      ],
    },
    {
      key: "3",
      icon: <WalletOutlined />,
      label: "Transactions",
      children: [
        {
          key: "5",
          label: (
            <Link href="http://localhost:3000/organization/overview/pending">
              Pending Contracts
            </Link>
          ),
        },
        { key: "6", label: "Option 6" },
      ],
    },
    {
      key: "4",
      icon: <UsergroupAddOutlined />,
      label: (
        <Link href="http://localhost:3000/organization//overview/employees/employeeList">
          Employees
        </Link>
      ),
    },
    {
      key: "sub1",

      label: (
        <Link href="http://localhost:3000/organization/overview/multiSig">
          Multisig
        </Link>
      ),
      icon: <KeyOutlined />,
      // children: [
      //   { key: "5", label: "Option 5" },
      //   { key: "6", label: "Option 6" },
      // ],
    },
  ];
  return (
    <div className="flex h-screen flex-col md:flex-row md:overflow-hidden justify-center pl-10  bg-slate-50">
      <div className="w-full flex-none md:w-24 pt-24">
        <Menu
          style={{
            borderRadius: 8,
            height: "228px",
            border: "solid 1px #1677FF",
          }}
          defaultSelectedKeys={["1"]}
          defaultOpenKeys={["sub1"]}
          mode="inline"
          theme="light"
          inlineCollapsed={collapsed}
          items={items}
        />
      </div>

      <div className="flex-grow  md:overflow-y-auto pt-2">
        <div className="flex  justify-end pt-2 mr-2">
          <BackBtn />
        </div>
        {children}
      </div>
    </div>
  );
}
