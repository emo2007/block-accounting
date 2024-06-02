export type LoginResponse = {
  token: string;
  token_expired_at: number;
  refresh_token: string;
  refresh_token_expired_at: number;
};
export type Organization = {
  id: string;
  name: string;
  address: string;
};
export type OrganizationsResponse = {
  items: Organization[];
};
export type NewOrgResponse = {
  id: string;
};
export type Ids = {
  id: string;
};
export type ParticipantsResponse = {
  ids: Ids[];
};
export type AddPartResponse = {
  name: string;
  position: string;
  wallet_address: string;
};
export type PublicKey = {
  public_key: string;
};
export type DeployResponse = {
  title: string;
  owners: PublicKey[];
  confirmations: string;
};
export type MultiSigResponse = {};

export type InvitationResponse = {
  token: string;
  token_expired_at: number;
  refresh_token: string;
  refresh_token_expired_at: number;
};
