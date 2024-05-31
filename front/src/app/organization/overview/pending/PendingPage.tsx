"use client";
import React, { useState, useEffect } from "react";
import { Steps, Typography, Card, Input, Space, Button, Select } from "antd";
import type { MenuProps } from "antd";
import { WalletOutlined } from "@ant-design/icons";
const { Title } = Typography;
type MenuItem = Required<MenuProps>["items"][number];
const data = [
  "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
  "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
  "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
  "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
  "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
];

export function PendingPage() {
  return (
    <div className="flex flex-col w-full h-full  gap-5 p-8  ">
      <Title style={{ color: "#302d43" }}>Pending Confirmations</Title>
      <Steps
        current={1}
        items={[
          {
            title: "Request",
          },
          {
            title: "Verify",
          },
          {
            title: "Success",
          },
        ]}
      />

      <Card className="flex flex-col  w-full ">
        <Space.Compact style={{ width: "100%" }}>
          <Input
            size="middle"
            defaultValue=""
            placeholder="Pay Out Transaction"
          />
          <Button style={{ color: "#4096ff", borderColor: "#4096ff" }}>
            Submit
          </Button>
        </Space.Compact>

        <div className="flex justify-end mt-5">
          <Button type="primary" style={{ width: "150px" }}>
            Execute
          </Button>
        </div>
      </Card>
      <Card className="flex flex-col  w-full ">
        <Space.Compact style={{ width: "100%" }}>
          <Input
            size="middle"
            defaultValue=""
            placeholder="Pay Out Transaction"
          />
          <Button style={{ color: "#4096ff", borderColor: "#4096ff" }}>
            Submit
          </Button>
        </Space.Compact>

        <div className="flex justify-end mt-5">
          <Button type="primary" style={{ width: "150px" }}>
            Execute
          </Button>
        </div>
      </Card>
      <Card className="flex flex-col  w-full ">
        <Space.Compact style={{ width: "100%" }}>
          <Input
            size="middle"
            defaultValue=""
            placeholder="Pay Out Transaction"
          />
          <Button style={{ color: "#4096ff", borderColor: "#4096ff" }}>
            Submit
          </Button>
        </Space.Compact>

        <div className="flex justify-end mt-5">
          <Button type="primary" style={{ width: "150px" }}>
            Execute
          </Button>
        </div>
      </Card>
    </div>
  );
}
