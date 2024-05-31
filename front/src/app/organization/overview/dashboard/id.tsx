"use client";
import React, { useState, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Card } from "antd";
import { UserOutlined } from "@ant-design/icons";
import { Button, List, Divider, Typography, Avatar, Skeleton } from "antd";

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
const count = 3;
const fakeDataUrl = `https://randomuser.me/api/?results=${count}&inc=name,gender,email,nat,picture&noinfo`;

const data = [
  "Ackee Blockchain is a team of auditors and white hat hackers who perform security audits and assessments for Ethereum and Solana.",
  "Global blockchain services company and Initial Coin Offering solutions provider",
  "AutoMinter is a decentralized no-code NFT collection generation platform.",
  "BANKEX will create smart contracts of any complexity for your projects in the Solidity language.",
  "Securing the DeFi ecosystem",
];

const { Title } = Typography;

export function OrgProfile() {
  const [initLoading, setInitLoading] = useState(true);
  const [loading, setLoading] = useState(false);
  const [dataEmployees, setData] = useState<DataType[]>([]);
  const [list, setList] = useState<DataType[]>([]);

  const router = useRouter();
  const pathname = useSearchParams();
  console.log(pathname.getAll("query"));

  const onNextPageHandler = () => {
    router.push("/organization/overview/employees");
  };
  const onMultisigPageHandler = () => {
    router.push("/organization/overview/multiSig");
  };
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
    <div className="flex flex-row w-full h-full bg-slate-50  p-8">
      <div className="flex flex-col w-11/12 ">
        <Title style={{ color: "#302d43", textIndent: 15 }}>Dashboard</Title>
        <Card
          title="Organization Name"
          bordered={false}
          style={{ width: "60%" }}
        >
          <p>Address</p>
          <p>Phone</p>
          <p>Description</p>
        </Card>
        <div className="flex  w-full justify-end ">
          <Button
            type="primary"
            size={"large"}
            style={{ width: "240px" }}
            onClick={onMultisigPageHandler}
          >
            Create Multisig contract
          </Button>
        </div>
        <Divider
          style={{ color: "#1677FF" }}
          orientation="left"
          orientationMargin="0"
        >
          <a href="#">Contracts</a>
        </Divider>
        <List
          bordered
          dataSource={data}
          renderItem={(item) => (
            <List.Item>
              <Typography.Text mark></Typography.Text> {item}
            </List.Item>
          )}
        />
        <div className="flex  w-full justify-end ">
          <Button
            type="primary"
            size={"large"}
            style={{ width: "180px", marginTop: 20 }}
            onClick={onNextPageHandler}
          >
            Add new employee
          </Button>
        </div>
        <Divider
          style={{ color: "#1677FF" }}
          orientation="left"
          orientationMargin="0"
        >
          <a href="http://localhost:3000/organization/employees/employeeList">
            Employee List
          </a>
        </Divider>
        <List
          className="demo-loadmore-list"
          loading={initLoading}
          itemLayout="horizontal"
          loadMore={loadMore}
          dataSource={list}
          renderItem={(item) => (
            <List.Item
              actions={[
                <a key="list-loadmore-edit">edit</a>,
                <a key="list-loadmore-more">more</a>,
              ]}
            >
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
    </div>
  );
}
