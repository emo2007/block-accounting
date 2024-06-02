// 3. Страница приглашения сотрудников (Employee Invitation Page)
// Генерация ссылки: Реализовать функционал для генерации и отправки пригласительной ссылки сотрудникам.
// Ввод SEED_KEY: Предоставить поле для ввода SEED_KEY при регистрации сотрудника, с возможностью генерации нового SEED_KEY для новых пользователей.

"use client";
import React, { useState, useEffect } from "react";
import { Input, Button, Typography, Space } from "antd";
import { MailOutlined, CopyOutlined } from "@ant-design/icons";
import { apiService } from "@/app/axios/global.service";

export function EmployeeInvitatonPage() {
  const { Text } = Typography;
  const [values, setValues] = useState({
    email: "ower@gmail.com",
    phone: "79999999999",
    telegram: "@ower",
  });

  async function sendInvite() {
    const data = await apiService.sentInvitation(
      "",
      "Alex",
      values,
      "short orient camp maple lend pole balance token pledge fat analyst badge art happy sadsad"
    );
    console.log(data);
  }

  return (
    <div className="flex relative  overflow-hidden flex-col w-2/3 h-3/4 items-center justify-center gap-10 bg-white border-solid border rounded-md border-neutral-300 p-10 text-neutral-500">
      <div className="w-full h-20    bg-[#1677FF] absolute top-0 flex items-center justify-center">
        <h1 className="text-white text-xl font-semibold">Invite</h1>
      </div>
      <div className="flex flex-col  w-9/12 pb-36 gap-5 items-start ">
        <Text type="secondary">Invite new Employee</Text>
        <Space.Compact style={{ width: "100%" }}>
          <Input
            size="large"
            name="email"
            addonBefore="E-mail"
            value={values.email}
            onChange={(e) =>
              setValues((prev) => ({ ...prev, email: e.target.value }))
            }
          />
          <Button
            size="large"
            icon={<MailOutlined />}
            style={{ width: 210 }}
            type="primary"
            onClick={sendInvite}
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
