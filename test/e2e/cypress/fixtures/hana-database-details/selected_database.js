export const selectedDatabase = {
  Id: 'fd44c254ccb14331e54015c720c7a1f2',
  Sid: 'HDD',
  Type: 'HANA Database',
  Hosts: [
    {
      Hostname: 'vmhdbdev01',
      Instance: '10',
      Features: 'HDB|HDB_WORKER',
      HttpPort: '51013',
      HttpsPort: '51014',
      StartPriority: '0.3',
      Status: 'SAPControl-GREEN',
      StatusBadge: 'badge-primary',
    },
    {
      Hostname: 'vmhdbdev02',
      Instance: '10',
      Features: 'HDB|HDB_WORKER',
      HttpPort: '51013',
      HttpsPort: '51014',
      StartPriority: '0.3',
      Status: 'SAPControl-GREEN',
      StatusBadge: 'badge-primary',
    },
  ],
};

export const attachedHosts = [
  {
    Name: 'vmhdbdev01',
    AgentId: '13e8c25c-3180-5a9a-95c8-51ec38e50cfc',
    Address: '10.100.1.11, 10.100.1.13',
    Provider: 'azure',
    Cluster: 'hana_cluster',
    ClusterId: '04b8f8c21f9fd8991224478e8c4362f8',
    Version: '0.7.1+git.dev42.1640084952.33229fc',
  },
  {
    Name: 'vmhdbdev02',
    AgentId: '0a055c90-4cb6-54ce-ac9c-ae3fedaf40d4',
    Address: '10.100.1.12',
    Provider: 'azure',
    Cluster: 'hana_cluster',
    ClusterId: '04b8f8c21f9fd8991224478e8c4362f8',
    Version: '0.7.1+git.dev42.1640084952.33229fc',
  },
];
