import { Injectable } from '@angular/core';
import {
  GoLoadServiceApi,
  Configuration,
  V1CreateAccountResponse,
  V1CreateSessionResponse,
  V1CreateDownloadTaskResponse,
  V1GetDownloadTaskListResponse,
  V1UpdateDownloadTaskResponse,
  V1CreateAccountRequest,
  V1CreateSessionRequest,
  V1CreateDownloadTaskRequest,
  V1DeleteDownloadTaskRequest,
  V1GetDownloadTaskFileRequest,
  V1GetDownloadTaskListRequest,
  V1UpdateDownloadTaskRequest,
  StreamResultOfV1GetDownloadTaskFileResponseFromJSON,
} from './api';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private readonly api: GoLoadServiceApi;

  constructor() {
    this.api = new GoLoadServiceApi(new Configuration({}));
  }

  public async createAccount(
    accountName: string,
    password: string
  ): Promise<V1CreateAccountResponse> {
    return this.api.goLoadServiceCreateAccount({
      body: {
        accountName: accountName,
        password: password,
      } as V1CreateAccountRequest,
    });
  }

  public async createSession(
    accountName: string,
    password: string
  ): Promise<V1CreateSessionResponse> {
    return this.api.goLoadServiceCreateSession({
      body: {
        accountName: accountName,
        password: password,
      } as V1CreateSessionRequest,
    });
  }

  public async createDownloadTask(
    accountId: string,
    downloadUrl: string
  ): Promise<V1CreateDownloadTaskResponse> {
    return this.api.goLoadServiceCreateDownloadTask({
      body: {
        accountId: accountId,
        downloadUrl: downloadUrl,
      } as V1CreateDownloadTaskRequest,
    });
  }

  public async deleteDownloadTask(
    accountId: string,
    taskId: string
  ): Promise<void> {
    await this.api.goLoadServiceDeleteDownloadTask({
      body: {
        accountId: accountId,
        taskId: taskId,
      } as V1DeleteDownloadTaskRequest,
    });
  }

  public async getDownloadTaskFile(
    accountId: string,
    taskId: string
  ): Promise<Observable<string>> {
    const response = await this.api.goLoadServiceGetDownloadTaskFileRaw({
      body: {
        accountId: accountId,
        taskId: taskId,
      } as V1GetDownloadTaskFileRequest,
    });

    return new Observable((subscriber) => {
      if (!response.raw.body) {
        subscriber.complete();
        return;
      }

      const bodyReader = response.raw.body.getReader();
      bodyReader.read().then(
        (rawBodyChunk) => {
          if (rawBodyChunk.done) {
            subscriber.complete();
            return;
          }

          const jsonBodyChunk =
            StreamResultOfV1GetDownloadTaskFileResponseFromJSON(
              rawBodyChunk.value
            );
          if (!jsonBodyChunk.result?.data) {
            return;
          }

          subscriber.next(btoa(jsonBodyChunk.result.data));
        },
        (error) => {
          subscriber.error(error);
        }
      );
    });
  }

  public async getDownloadTaskList(
    accountId: string,
    offset: number,
    limit: number
  ): Promise<V1GetDownloadTaskListResponse> {
    return this.api.goLoadServiceGetDownloadTaskList({
      body: {
        accountId: accountId,
        offset: `${offset}`,
        limit: `${limit}`,
      } as V1GetDownloadTaskListRequest,
    });
  }

  public async updateDownloadTask(
    accountId: string,
    taskId: string,
    status: string
  ): Promise<V1UpdateDownloadTaskResponse> {
    return this.api.goLoadServiceUpdateDownloadTask({
      body: {
        accountId: accountId,
        taskId: taskId,
        status: status,
      } as V1UpdateDownloadTaskRequest,
    });
  }
}
