"use client";
import React, { useState, useEffect } from "react";
import { Button, Menu, List, Typography, Avatar, Skeleton } from "antd";
import { UserOutlined } from "@ant-design/icons";
import type { MenuProps } from "antd";
import { PayOutBtn } from "@/app/ui/PayOutBtn";
import { apiService } from "@/app/axios/global.service";
import useOrganizationsHooks from "@/hooks/organizations";
import { useSearchParams } from "next/navigation";
import { Organization } from "@/app/axios/api-types";
const count = 8;
const fakeDataUrl = `https://randomuser.me/api/?results=${count}&inc=name,gender,email,nat,picture&noinfo`;
interface DataType {
  gender?: string;
  name: {
    title?: string;
    first?: string;
    last?: string;
  };
  email?: string;
  picture: {
    large?: string;
    medium?: string;
    thumbnail?: string;
  };
  nat?: string;
  loading: boolean;
}
const { Title } = Typography;
type MenuItem = Required<MenuProps>["items"][number];
export function EmployeeList() {
  const {
    organizations,
    filteredOrganization,
    setOrganizations,
    loadOrganizations,
  } = useOrganizationsHooks();
  const [collapsed, setCollapsed] = useState(true);
  const [initLoading, setInitLoading] = useState(true);
  const [loading, setLoading] = useState(false);
  const [dataEmployees, setData] = useState<DataType[]>([]);
  const [list, setList] = useState<DataType[]>([]);
  const [organization, setOrganization] = useState<Organization>({
    id: "",
    name: "",
    address: "",
  });
  const loadEmployees = async (id: string) => {
    const data: any = await apiService.getEmployees(id, []);
    setList(data.data.participants);
  };
  const searchParams = useSearchParams();
  useEffect(() => {
    const id = searchParams.get("id") || "";
    if (filteredOrganization) {
      loadEmployees(id);
      setOrganization(filteredOrganization);
      setInitLoading(false);
    }
  }, [organizations]);
  useEffect(() => {
    fetch(fakeDataUrl)
      .then((res) => res.json())
      .then((res) => {
        setInitLoading(false);
        setData(res.results);
        setList(res.results);
      });
  }, []);
  const onLoadMore = () => {
    setLoading(true);
    loadEmployees(organization.id);
  };
  const loadMore =
    !initLoading && !loading ? (
      <div
        style={{
          textAlign: "center",
          marginTop: 12,
          height: 32,
          lineHeight: "32px",
        }}
      >
        <Button onClick={onLoadMore}>loading more</Button>
      </div>
    ) : null;
  return (
    <div className="flex flex-col w-full h-full  gap-5 pb-20 px-30 p-10">
      <Title style={{ color: "#302d43", textIndent: 15 }}>Employee List</Title>
      <List
        className="demo-loadmore-list"
        loading={initLoading}
        itemLayout="horizontal"
        loadMore={loadMore}
        dataSource={list}
        renderItem={(item) => {
          console.log(item);

          return (
            <List.Item actions={[<PayOutBtn />]}>
              <Skeleton avatar title={false} loading={item.loading} active>
                <List.Item.Meta
                  avatar={<Avatar icon={<UserOutlined />} />}
                  title={<a href="https://ant.design">{item.name}</a>}
                  description="1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71"
                />
                <div>wallet address</div>
              </Skeleton>
            </List.Item>
          );
        }}
      />
    </div>
  );
}
