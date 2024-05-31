"use client";
import React, { useState, useEffect } from "react";
import { Button, Menu, List, Typography, Avatar, Skeleton } from "antd";
import { UserOutlined } from "@ant-design/icons";
import type { MenuProps } from "antd";
import { PayOutBtn } from "@/app/ui/PayOutBtn";
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
  const [collapsed, setCollapsed] = useState(true);
  const [initLoading, setInitLoading] = useState(true);
  const [loading, setLoading] = useState(false);
  const [dataEmployees, setData] = useState<DataType[]>([]);
  const [list, setList] = useState<DataType[]>([]);
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
    setList(
      dataEmployees.concat(
        [...new Array(count)].map(() => ({
          loading: true,
          name: {},
          picture: {},
        }))
      )
    );
    fetch(fakeDataUrl)
      .then((res) => res.json())
      .then((res) => {
        const newData = dataEmployees.concat(res.results);
        setData(newData);
        setList(newData);
        setLoading(false);
        // Resetting window's offsetTop so as to display react-virtualized demo underfloor.
        // In real scene, you can using public method of react-virtualized:
        // https://stackoverflow.com/questions/46700726/how-to-use-public-method-updateposition-of-react-virtualized
        window.dispatchEvent(new Event("resize"));
      });
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
    <div className="flex flex-col w-full h-full  gap-5 pb-20 px-30 p-8">
      <Title style={{ color: "#302d43", textIndent: 15 }}>Employee List</Title>
      <List
        className="demo-loadmore-list"
        loading={initLoading}
        itemLayout="horizontal"
        loadMore={loadMore}
        dataSource={list}
        renderItem={(item) => (
          <List.Item actions={[<PayOutBtn />]}>
            <Skeleton avatar title={false} loading={item.loading} active>
              <List.Item.Meta
                avatar={<Avatar icon={<UserOutlined />} />}
                title={<a href="https://ant.design">{item.name?.last}</a>}
                description="1Lbcfr7sAHTD9CgdQo3HTMTkV8LK4ZnX71"
              />
              <div>wallet address</div>
            </Skeleton>
          </List.Item>
        )}
      />
    </div>
  );
}
