interface TreeNodeType {
  pk: string;
  invited_by: string;
  children?: TreeNodeType[];
}

export function buildHierarchy(
  data: TreeNodeType[],
  parent: TreeNodeType | null = null,
): TreeNodeType[] {
  const children = data.filter(
    (item) => item.invited_by === (parent ? parent.pk : null),
  );
  return children.map((child) => ({
    ...child,
    children: buildHierarchy(data, child),
  }));
}
