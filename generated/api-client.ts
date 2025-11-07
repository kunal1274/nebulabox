// Auto-generated from schema/nebulabox.schema.json
// DO NOT EDIT MANUALLY

import { ApiClient } from "../web/dashboard/src/lib/api";
import type { Container, Workspace } from "./types";

export class NebulaBoxAPI {
  constructor(private client: ApiClient) {}

  async containersList(
    query?: { all?: boolean; workspaceId?: string,null }
  ): Promise<any> {
    return this.client.request('GET', `/api/containers`, query);
  }

  async containersGet(
    id: string
  ): Promise<Container> {
    return this.client.request('GET', `/api/containers/${id}`);
  }

  async containersCreate(
    data: any
  ): Promise<Container> {
    return this.client.request('POST', `/api/containers/run`, data);
  }

  async containersStart(
    id: string
  ): Promise<any> {
    return this.client.request('POST', `/api/containers/${id}/start`);
  }

  async containersStop(
    id: string
  ): Promise<any> {
    return this.client.request('POST', `/api/containers/${id}/stop`);
  }

  async containersDelete(
    id: string
  ): Promise<any> {
    return this.client.request('DELETE', `/api/containers/${id}`);
  }

  async imagesList(
    
  ): Promise<any> {
    return this.client.request('GET', `/api/images`);
  }

  async imagesPull(
    data: any
  ): Promise<any> {
    return this.client.request('POST', `/api/images/pull`, data);
  }

  async imagesBuild(
    data: any
  ): Promise<any> {
    return this.client.request('POST', `/api/images/build`, data);
  }

  async workspacesList(
    
  ): Promise<any> {
    return this.client.request('GET', `/api/workspaces`);
  }

  async workspacesGet(
    id: string
  ): Promise<Workspace> {
    return this.client.request('GET', `/api/workspaces/${id}`);
  }

  async workspacesCreate(
    data: any
  ): Promise<Workspace> {
    return this.client.request('POST', `/api/workspaces`, data);
  }

}