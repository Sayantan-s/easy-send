import { v4 as uuid } from "uuid";
import {
  IApiFailureResponsePayload,
  IApiFailureResponseTemplate,
  IApiSuccessResponsePayload,
  IApiSuccessResponseTemplate,
} from "./types";

class ApiResponse {
  static success<TData>({
    status,
    data,
    message,
  }: IApiSuccessResponsePayload<TData>) {
    const responseTemplate: IApiSuccessResponseTemplate<TData> = {
      status,
      data,
      requestId: uuid(),
      message: message || "Operation Successfull",
    };
    return new Response(JSON.stringify(responseTemplate), { status });
  }

  static failure<TError>({
    status,
    error,
    message,
  }: IApiFailureResponsePayload<TError>) {
    const responseTemplate: IApiFailureResponseTemplate<TError> = {
      status,
      error,
      requestId: uuid(),
      message: message || "Operation Failure",
    };
    return new Response(JSON.stringify(responseTemplate), { status });
  }
}

export default class Communication {
  static response = ApiResponse;
}
