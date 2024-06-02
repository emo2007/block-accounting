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
  InvitationResponse,
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
    organizationId: string,
    ids: string[]
  ): Promise<AxiosResponse<ParticipantsResponse>> {
    return await globalService.post(
      `organizations/${organizationId}/participants/fetch`,
      {
        data: {
          ids: [...ids, organizationId],
        },
      }
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
    organizationId: string,
    title: string,
    owners: string[],
    confirmations: number
  ): Promise<AxiosResponse<DeployResponse>> {
    return await globalService.post(
      `organizations/${organizationId}/multisig`,
      {
        title,
        owners,
        confirmations,
      }
    );
  }

  async getAllMultisigsByOrganizationId(
    organizationId: string
  ): Promise<AxiosResponse<MultiSigResponse>> {
    return await globalService.get(`organizations/${organizationId}/multisig`);
  }

  async sentInvitation(
    hash: string,
    name: string,
    credentals: {
      email: string;
      phone: string;
      telegram: string;
    },
    mnemonic: string
  ): Promise<AxiosResponse<InvitationResponse>> {
    return await globalService.post(`/invite/${hash}/join`, {
      data: {
        name,
        credentals,
        mnemonic,
      },
    });
  }
}

export const apiService = new AccountingService();
