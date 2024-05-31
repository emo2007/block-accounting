"use client";
import React, { useState } from "react";
import { Button, Form, Input, Select, Space } from "antd";
import { useRouter } from "next/navigation";
export function OrgForm({ setFormData }) {
  const { Option } = Select;
  const prefixSelector = (
    <Form.Item name="prefix" noStyle>
      <Select style={{ width: 70 }}>
        <Option value="86">+86</Option>
        <Option value="87">+87</Option>
      </Select>
    </Form.Item>
  );
  const router = useRouter();
  const onNextPageHandler = () => {
    router.push("/organization/dashboard");
  };
  return (
    <div className="bg-white p-10 flex items-center justify-start pl-36">
      <Form
        name="complex-form"
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 16 }}
        style={{ maxWidth: 600 }}
      >
        <div>
          <Form.Item label="Organization name">
            <Space>
              <Form.Item
                noStyle
                rules={[
                  {
                    required: true,
                    message: "Organization name is required",
                  },
                ]}
              >
                <Input
                  name="name"
                  style={{ width: 350 }}
                  onInput={(element: any) =>
                    setFormData((prev: object) => ({
                      ...prev,
                      [element.target.name]: element.target.value,
                    }))
                  }
                />
              </Form.Item>
            </Space>
          </Form.Item>
          <Form.Item label="Address">
            <Space>
              <Form.Item
                name="address"
                noStyle
                rules={[{ required: true, message: "Address is required" }]}
              >
                <Input
                  name="address"
                  onInput={(element: any) =>
                    setFormData((prev: object) => ({
                      ...prev,
                      [element.target.name]: element.target.value,
                    }))
                  }
                  style={{ width: 350 }}
                />
              </Form.Item>
            </Space>
          </Form.Item>
          {/*<Form.Item label="Phone Number">*/}
          {/*  <Space>*/}
          {/*    <Form.Item*/}
          {/*      noStyle*/}
          {/*      rules={[*/}
          {/*        {*/}
          {/*          required: true,*/}
          {/*          message: "Please input phone number",*/}
          {/*        },*/}
          {/*      ]}*/}
          {/*    >*/}
          {/*      <Input*/}
          {/*        name="phone"*/}
          {/*        addonBefore={prefixSelector}*/}
          {/*        style={{ width: 350 }}*/}
          {/*        onInput={(e: any) =>*/}
          {/*          setFormData((prev: object) => ({*/}
          {/*            ...prev,*/}
          {/*            [e.target.name]: e.target.value,*/}
          {/*          }))*/}
          {/*        }*/}
          {/*        type="number"*/}
          {/*        required*/}
          {/*        minLength={9}*/}
          {/*        maxLength={9}*/}
          {/*      />*/}
          {/*    </Form.Item>*/}
          {/*  </Space>*/}
          {/*</Form.Item>*/}
        </div>
      </Form>
    </div>
  );
}
