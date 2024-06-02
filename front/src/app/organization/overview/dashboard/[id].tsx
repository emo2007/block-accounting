"use client";
import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useSearchParams } from "next/navigation";
import { UserOutlined } from "@ant-design/icons";
import { Button, List, Divider, Typography, Avatar, Skeleton } from "antd";
import { OrganizationCard } from "@/app/orgCreate/OrgCard";
import { Organization } from "@/app/axios/api-types";
import useOrganizationsHooks from "@/hooks/organizations";
import { apiService } from "@/app/axios/global.service";

type OrgItemProps = {
  element: Organization;
};
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
  const { organizations, setOrganizations, loadOrganizations } =
    useOrganizationsHooks();
  const [organization, setOrganization] = useState<Organization>({
    id: "",
    name: "",
    address: "",
  });
  const [initLoading, setInitLoading] = useState(true);
  const [loading, setLoading] = useState(false);
  const [dataEmployees, setData] = useState<DataType[]>([]);
  const [list, setList] = useState<any[]>([]);

  const router = useRouter();
  const searchParams = useSearchParams();
  useEffect(() => {
    const id = searchParams.get("id") || "";
    const filteredOrganization = organizations.find(
      (element) => element.id === id
    );

    if (filteredOrganization) {
      loadEmployees(id);
      setOrganization(filteredOrganization);
      setInitLoading(false);
    }
  }, [organizations]);

  useEffect(() => {
    loadOrganizations();
  }, []);

  const loadEmployees = async (id: string) => {
    const data: any = await apiService.getEmployees(id, []);
    setList(data.data.participants);
  };

  const onNextPageHandler = () => {
    router.push("/organization/overview/employees");
  };
  const onMultisigPageHandler = () => {
    router.push("/organization/overview/multiSig/?id=" + organization.id);
  };
  const onLoadMore = () => {
    loadEmployees(organization.id);
    setInitLoading(false);
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

        <div style={{ width: "60%" }} className="flex flex-col  ">
          {organization && <OrganizationCard element={organization} />}
        </div>

        {/* <Card
          title="Organization Name"
          bordered={false}
          style={{ width: "60%" }}
        >
          <p>Address</p>
          <p>Phone</p>
          <p>Description</p>
        </Card> */}
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
          renderItem={(item) => {
            console.log(item);

            return (
              <List.Item
                actions={[
                  <a key="list-loadmore-edit">edit</a>,
                  <a key="list-loadmore-more">more</a>,
                ]}
              >
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
    </div>
  );
}
