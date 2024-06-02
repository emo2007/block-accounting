// 2. Страница создания организации (Organization Creation Page)
// Поля для ввода: Предоставить поля для ввода названия организации и дополнительной информации, такой как адрес и контактные данные.
// Кнопка создания: Иметь кнопку "Создать", которая активируется только после заполнения всех необходимых полей.
// Обратная связь: Показывать сообщения об ошибках или подтверждение успешного создания организации.

//* <h1>{seed.join("\n")}</h1> */shtuchka kak map

"use client";
import React, { useEffect } from "react";
import { Button, Modal } from "antd";
import { useState } from "react";
import { apiService } from "../axios/global.service";
import { OrgForm } from "./OrgForm";
import { OrganizationCard } from "./OrgCard";
import { FolderOpenTwoTone } from "@ant-design/icons";
import useOrganizationsHooks from "@/hooks/organizations";
import {
  Organization,
  NewOrgResponse,
  OrganizationsResponse,
} from "../axios/api-types";
import Cookies from "js-cookie";
export function OrgCreatePage() {
  //const [organizations, setOrganizations] = useState<Organization[]>([]);
  const { organizations, setOrganizations } = useOrganizationsHooks();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    address: "",
  });
  const onFinish = async (values: any) => {
    handleOk();
    setOrganizations((prev: any) => [...prev, formData]);

    // setFormData({});
    const result = await apiService.newOrganization(
      formData.name,
      formData.address
    );
    if (result) {
      loadOrganizations();
    }
  };
  const showModal = () => {
    setIsModalOpen(true);
  };

  const handleOk = () => {
    setIsModalOpen(false);
  };

  console.log(organizations);

  const handleCancel = () => {
    setIsModalOpen(false);
  };
  useEffect(() => {
    loadOrganizations();
  }, []);

  const loadOrganizations = async () => {
    const result = await apiService.getOrganizations();
    if (result) {
      setOrganizations(result.data.items || []);
    }
  };
  return (
    <>
      <div className="flex relative  overflow-hidden  flex-col w-2/3 h-3/4 items-center justify-center z-30 gap-5 bg-white border-solid border rounded-md border-neutral-300 text-neutral-500">
        <div className="w-full h-20   bg-[#1677FF] absolute top-0 flex items-center z-40 justify-center">
          <h1 className="text-white text-xl font-semibold">
            Your Organizations
          </h1>
        </div>
        <div></div>
        <div className="flex flex-col relative  w-full h-3/4  items-center  overflow-scroll gap-10  p-10 z-0">
          {organizations.length && (
            <FolderOpenTwoTone style={{ fontSize: "400%" }} />
          )}
          {organizations.length ? (
            organizations.map((element: any) => {
              return (
                <div className="flex flex-col min-w-full  ">
                  <OrganizationCard element={element} />
                </div>
              );
            })
          ) : (
            <em>Your Organization list is currently empty.</em>
          )}

          <div>
            <Modal
              width={1000}
              centered
              title="Please input information"
              open={isModalOpen}
              okText="Submit"
              onOk={onFinish}
              okButtonProps={{
                disabled: !(Object.values(formData).length === 2),
              }}
              onCancel={handleCancel}
            >
              <OrgForm setFormData={setFormData} />
            </Modal>
          </div>
        </div>
        <div className=" flex z-40 ">
          <Button
            style={{ width: "150px" }}
            size="large"
            type="primary"
            onClick={showModal}
          >
            Create
          </Button>
        </div>
      </div>
    </>
  );
}
