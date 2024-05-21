"use client";
import React, { useState } from "react";
import { Input } from "antd";
import { Typography } from "antd";
import { Card } from "antd";
import type { InputNumberProps } from "antd";
import { Col, InputNumber, Row, Slider, Space, Button } from "antd";

const { Title } = Typography;

export function MultisigPage() {
  const [inputValue, setInputValue] = useState(1);
  const onChange: InputNumberProps["onChange"] = (newValue) => {
    setInputValue(newValue as number);
  };
  return (
    <div className="flex flex-col w-full h-full px-28 py-20 gap-10 ">
      <div className="flex flex-col w-1/3">
        <Title level={3}>Create a new Multisig</Title>
        <Input size="large" placeholder="Multisig Name/Label" />
      </div>
      <div className="flex  w-full  ">
        <Card style={{ width: "100%" }}>
          <Title level={4}>Signers</Title>
          <div className="flex flex-row gap-5">
            <div className="flex flex-col gap-2 w-1/4">
              <Input placeholder="Name" />
              <Input placeholder="Name" />
              <Input placeholder="Name" />
            </div>
            <div className="flex flex-col gap-2 w-full">
              <Input placeholder="Pubic Key" />
              <Input placeholder="Pubic Key" />
              <Input placeholder="Pubic Key" />
            </div>
          </div>
          <div className="flex  w-full justify-end mt-5">
            <Button size={"large"} type="primary">
              Add Signer
            </Button>
          </div>
        </Card>
      </div>
      <div className="mt-20">
        <Card style={{ width: "100%" }}>
          <Row
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
            }}
          >
            <Col span={12}>
              <Slider
                min={1}
                max={20}
                onChange={onChange}
                value={typeof inputValue === "number" ? inputValue : 0}
              />
            </Col>
            <Col span={4}>
              <InputNumber
                min={1}
                max={20}
                style={{ margin: "0 16px" }}
                value={inputValue}
                onChange={onChange}
              />
            </Col>
          </Row>
        </Card>
      </div>
      <div className="flex  w-full justify-end pr-5">
        <Button size={"large"} type="primary">
          Create Multisig
        </Button>
      </div>
    </div>
  );
}
