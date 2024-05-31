"use client";
import React, { useState, useEffect } from "react";
import { Select, Typography, Card, Divider, Button } from "antd";
import { WalletOutlined } from "@ant-design/icons";
const handleChange = (value: string) => {
  console.log(`selected ${value}`);
};

const { Title } = Typography;

export function AgreementPage() {
  return (
    <div className="flex flex-col w-full p-8 ">
      <Title style={{ color: "#302d43" }}>Agreement</Title>

      <Card className="w-full ">
        <Divider
          style={{ color: "#1677FF" }}
          orientation="left"
          orientationMargin="0"
        >
          <a href="#">Owners</a>
        </Divider>
        <div className="flex flex-col gap-2">
          <Select
            suffixIcon={<WalletOutlined />}
            defaultValue={""}
            style={{ width: "full" }}
            allowClear
            options={[{ value: "", label: "" }]}
          />
          <Select
            suffixIcon={<WalletOutlined />}
            defaultValue={""}
            style={{ width: "full" }}
            allowClear
            options={[{ value: "", label: "" }]}
          />
          <Select
            suffixIcon={<WalletOutlined />}
            defaultValue={""}
            style={{ width: "full" }}
            allowClear
            options={[{ value: "", label: "" }]}
          />
          <Select
            suffixIcon={<WalletOutlined />}
            defaultValue={""}
            style={{ width: "full" }}
            allowClear
            options={[{ value: "", label: "" }]}
          />
          <Select
            suffixIcon={<WalletOutlined />}
            defaultValue={""}
            style={{ width: "full" }}
            allowClear
            options={[{ value: "", label: "" }]}
          />
        </div>
      </Card>
      <div className="flex  w-full justify-end mt-8">
        <Button size={"large"} style={{ width: 150 }} type="primary">
          Confirm
        </Button>
      </div>
    </div>
  );
}
