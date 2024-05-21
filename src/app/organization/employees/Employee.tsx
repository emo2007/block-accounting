// 3. Страница приглашения сотрудников (Employee Invitation Page)
// Генерация ссылки: Реализовать функционал для генерации и отправки пригласительной ссылки сотрудникам.
// Ввод SEED_KEY: Предоставить поле для ввода SEED_KEY при регистрации сотрудника, с возможностью генерации нового SEED_KEY для новых пользователей.

"use client";
import React, { useState, useEffect } from "react";
import { Input, Button, Typography, Space } from "antd";
import { MailOutlined, CopyOutlined } from "@ant-design/icons";
export function EmployeePage() {
  const { Text } = Typography;
  return (
    <div className="flex relative  overflow-hidden flex-col w-2/3 h-3/4 items-center justify-center gap-10 bg-white border-solid border rounded-md border-neutral-300 text-neutral-500">
      <div className="w-full h-20    bg-[#1677FF] absolute top-0 flex items-center justify-center">
        <h1 className="text-white text-xl font-semibold">Invite</h1>
      </div>
      <div className="flex flex-col  w-9/12 pb-36 gap-5 items-start ">
        <Text type="secondary">Invite new Employee</Text>
        <Space.Compact style={{ width: "100%" }}>
          <Input size="large" addonBefore="E-mail" />
          <Button
            size="large"
            icon={<MailOutlined />}
            style={{ width: 210 }}
            type="primary"
          >
            Send Invite
          </Button>
        </Space.Compact>
        <Space.Compact style={{ width: "100%" }}>
          <Input size="large" placeholder="invitation link" />
          <Button
            size="large"
            icon={<CopyOutlined />}
            style={{ width: 210 }}
            type="primary"
          >
            Copy Link Invite
          </Button>
        </Space.Compact>
      </div>
    </div>
  );
}
