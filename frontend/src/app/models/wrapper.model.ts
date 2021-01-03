export interface Wrapper<T> {
  data: T[];
  totalCount: number;
  pageSize: number;
  page: number;
}
