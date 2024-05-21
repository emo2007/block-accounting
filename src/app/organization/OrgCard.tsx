"use client";
import React from "react";

import { Card } from "antd";
import { useState, useEffect, FC } from "react";
import { useRouter } from "next/navigation";
type OrgData = {
  name: string;
  address: string;
  phone: number;
};
type OrgItemProps = {
  element: OrgData;
};

export const OrganizationCard: FC<OrgItemProps> = ({ element }) => {
  const router = useRouter();
  const onNextPageHandler = () => {
    router.push("/organization/dashboard");
  };
  return (
    <>
      <Card
        title={element.name}
        type="inner"
        extra={
          <a onClick={onNextPageHandler} href="#">
            More
          </a>
        }
      >
        <p>{element.address}</p>
        <p>{element.phone}</p>
      </Card>
    </>
  );
};
