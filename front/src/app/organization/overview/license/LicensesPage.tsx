"use client";
import React, { useState, useEffect } from "react";
import { Typography, Button, Card, Input } from "antd";
import { WalletOutlined } from "@ant-design/icons";

const { Title } = Typography;
const data = [
  "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
  "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
  "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
  "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
  "1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71",
];

export function LicensesPage() {
  return (
    <div className="flex flex-col w-full p-8  ">
      <Title style={{ color: "#302d43" }}>Licenses</Title>

      <Card style={{ width: "100%" }}>
        <div className=" flex flex-row w-full gap-10">
          <div className="flex flex-col gap-2 w-1/4">
            <Title level={4}>Owners</Title>
            <Input placeholder="Name" />
            <Input placeholder="Name" />
            <Input placeholder="Name" />
            <Input placeholder="Name" />
            <Input placeholder="Name" />
          </div>
          <div className="flex flex-col gap-2 w-full">
            <Title level={4}>Shares</Title>

            <Input placeholder="Input information" />
            <Input placeholder="Input information" />
            <Input placeholder="Input information" />
            <Input placeholder="Input information" />
            <Input placeholder="Input information" />
          </div>
        </div>
      </Card>
      <div className="flex  w-full justify-end mt-5">
        <Button size={"large"} style={{ width: "150px" }} type="primary">
          Confirm
        </Button>
      </div>
    </div>
  );
}
//suffixIcon={<WalletOutlined />}
