"use client";
import { useState, useEffect, FC } from "react";

import { Card } from "antd";
import React from "react";
// type SeedData = {
//   seed: string;
// };
type SeedItemProps = {
  seed: string;
};
export const SeedItem: FC<SeedItemProps> = ({ seed }) => {
  console.log(seed);

  return (
    <Card size="small" style={{ minWidth: 100 }}>
      <p>{seed}</p>
    </Card>
  );
};
