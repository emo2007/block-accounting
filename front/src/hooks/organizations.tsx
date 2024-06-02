import { useState } from "react";
import { Organization } from "@/app/axios/api-types";
import { apiService } from "@/app/axios/global.service";
import { useSearchParams } from "next/navigation";

export default function useOrganizationsHooks() {
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [filteredOrganization, setFilteredOrganization] = useState<
    Organization[]
  >([]);
  const loadOrganizations = async () => {
    const result = await apiService.getOrganizations();
    if (result) {
      setOrganizations(result.data.items || []);
    }
    const searchParams = useSearchParams();
    const id = searchParams.get("id") || "";
    const filteredOrg = organizations.find((element) => element.id === id);

    if (filteredOrg) {
      setFilteredOrganization(filteredOrg || []);
    }
  };

  return {
    organizations,
    filteredOrganization,
    setOrganizations,
    loadOrganizations,
  };
}
