export type Filter<T> = Partial<Record<keyof T, any>>;

export type FindOptions<T> = {
  projection?: Partial<Record<keyof T, 0 | 1>>;
  sort?: Partial<Record<keyof T, 1 | -1>>;
  limit?: number;
  skip?: number;
};

export interface DataLayer<T> {
  create(data: Partial<T>): Promise<T>;
  find(filter?: Filter<T>, options?: FindOptions<T>): Promise<T[]>;
  findById(id: string | number): Promise<T | null>;
  update(id: string | number, data: Partial<T>): Promise<T | null>;
  delete(id: string | number): Promise<boolean>;
  paginate(
    filter?: Filter<T>,
    page?: number,
    perPage?: number,
    options?: FindOptions<T>
  ): Promise<{ items: T[]; total: number; page: number; perPage: number }>;
  model?: T;
}
