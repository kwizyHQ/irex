import { DataLayer, Filter, FindOptions } from "./dl.types";
{{ range .Models -}}
import {{ title .}}Model, { {{ title .}}Schema } from "./{{ lower . }}"
{{ end }}
export function mongooseAdapter<T>(model: any): DataLayer<T> {
  return {
    async create(data) {
      const doc = await model.create(data);
      return doc.toObject();
    },

    async find(filter = {}, options = {}) {
      let q = model.find(filter);
      if (options.projection) q = q.select(options.projection);
      if (options.sort) q = q.sort(options.sort);
      if (options.limit) q = q.limit(options.limit);
      if (options.skip) q = q.skip(options.skip);
      return (await q.exec()).map((d: any) => d.toObject());
    },

    async findById(id) {
      const doc = await model.findById(id).exec();
      return doc ? doc.toObject() : null;
    },

    async update(id, data) {
      const doc = await model.findByIdAndUpdate(id, data, { new: true }).exec();
      return doc ? doc.toObject() : null;
    },

    async delete(id) {
      return !!(await model.findByIdAndDelete(id).exec());
    },

    async paginate(filter = {}, page = 1, perPage = 10) {
      const offset = (page - 1) * perPage;
      const [items, total] = await Promise.all([
        model.find(filter).skip(offset).limit(perPage).exec(),
        model.countDocuments(filter).exec(),
      ]);
      return {
        items: items.map((d: any) => d.toObject()),
        total,
        page,
        perPage,
      };
    },

    model,
  };
}

const DL = {
{{ range .Models -}}
    {{ title . }}Model : mongooseAdapter<typeof {{ title .}}Schema>({{ title . }}Model),
{{ end }}
}

export default DL