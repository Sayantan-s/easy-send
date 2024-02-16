export interface IApiSuccessResponseTemplate<TData> {
  status: 200 | 201;
  message?: "Operation Successfull";
  data: TData;
  requestId: string;
}

export type IApiSuccessResponsePayload<TData> = Pick<
  IApiSuccessResponseTemplate<TData>,
  "status" | "data" | "message"
>;

export interface IApiFailureResponseTemplate<TError> {
  status: 400 | 401 | 402 | 403 | 404;
  message?: "Operation Failure";
  error: TError;
  requestId: string;
}

export type IApiFailureResponsePayload<TError> = Pick<
  IApiFailureResponseTemplate<TError>,
  "status" | "error" | "message"
>;
