import { createPromiseClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { UserService } from '../types/proto/api/v1/user_service_pb';
import { RegisterUserRequest, RegisterUserResponse } from '../types/proto/api/v1/user_service_pb';

const transport = createConnectTransport({
  baseUrl: 'http://localhost:8081',
});

export const userClient = createPromiseClient(UserService, transport);

export async function registerUser(
  username: string,
  email: string,
  password: string,
  nickname: string,
  phone: string
): Promise<RegisterUserResponse> {
  const request = new RegisterUserRequest({
    username,
    email,
    password,
    nickname,
    phone,
  });

  return await userClient.registerUser(request);
}
