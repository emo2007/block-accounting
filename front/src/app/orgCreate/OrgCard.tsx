"use client";
import React from "react";

import { Card } from "antd";
import { useState, useEffect, FC } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Organization } from "../axios/api-types";

type OrgItemProps = {
  element: Organization;
};

export const OrganizationCard: FC<OrgItemProps> = ({ element }) => {
  const router = useRouter();
  const id: any = element.id;

  return (
    <>
      <Card
        title={element.name}
        type="inner"
        extra={
          <Link
            href={{ pathname: "/organization/overview/dashboard/", query: id }}
          >
            More
          </Link>
        }
      >
        <p>{element.address}</p>
        <p>{element.name}</p>
      </Card>
    </>
  );
};
