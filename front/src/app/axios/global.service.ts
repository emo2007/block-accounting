import { globalService } from "../axios/global-api";
import { AxiosResponse } from "axios";
import Cookies from "js-cookie";
import {
  LoginResponse,
  OrganizationsResponse,
  NewOrgResponse,
  ParticipantsResponse,
  AddPartResponse,
  DeployResponse,
  MultiSigResponse,
} from "./api-types";
export class AccountingService {
  async login(seedKey: string): Promise<AxiosResponse<LoginResponse>> {
    console.log(seedKey);

    return await globalService.post("login", {
      mnemonic: seedKey,
    });
  }

  async register(seedKey: string) {
    return await globalService.post("join", {
      mnemonic: seedKey,
    });
  }

  async newOrganization(
    name: string,
    address: string
  ): Promise<AxiosResponse<NewOrgResponse>> {
    return await globalService.post("organizations", {
      name,
      address,
    });
  }

  async getOrganizations(): Promise<AxiosResponse<OrganizationsResponse>> {
    return await globalService.post("organizations/fetch", {
      data: {
        limit: 5,
      },
    });
  }

  async getEmployees(
    organizationId: string
  ): Promise<AxiosResponse<ParticipantsResponse>> {
    return await globalService.get(
      `organizations/${organizationId}/participants`
    );
  }

  async addEmployee(
    name: string,
    position: string,
    wallet_address: string,
    organizationId: string
  ): Promise<AxiosResponse<AddPartResponse>> {
    return await globalService.post(
      `organizations/${organizationId}/participants`,
      {
        name,
        position,
        wallet_address,
      }
    );
  }

  // POST /organizations/{organization_id}/multisig
  async deployMultisig(
    organization_id: string,
    title: string,
    owners: string[],
    confirmations: number
  ): Promise<AxiosResponse<DeployResponse>> {
    return await globalService.post(`${organization_id}/multi-sig`, {
      title,
      owners,
      confirmations,
    });
  }

  async getAllMultisigsByOrganizationId(
    organizationId: string
  ): Promise<AxiosResponse<MultiSigResponse>> {
    return await globalService.get(`organizations/${organizationId}/multisig`);
  }
}

export const apiService = new AccountingService();
