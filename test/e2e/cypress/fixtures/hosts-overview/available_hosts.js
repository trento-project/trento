const availableHosts = [
  ['a09d9cf3-46c1-505c-8fb8-4b0a71a9114e', 'vmdrbdprd01'],
  ['927901fa-2c87-524e-b18c-3ef5187f504f', 'vmdrbdprd02'],
  ['116d49bd-85e1-5e59-b820-83f66db8800c', 'vmnwprd01'],
  ['4b30a6af-4b52-5bda-bccb-f2248a12c992', 'vmnwprd02'],
  ['a3297d85-5e8b-5ac5-b8a3-55eebc2b8d12', 'vmnwprd03'],
  ['0fc07435-7ee2-54ca-b0de-fb27ffdc5deb', 'vmnwprd04'],
  ['9cd46919-5f19-59aa-993e-cf3736c71053', 'vmhdbprd01'],
  ['b767b3e9-e802-587e-a442-541d093b86b9', 'vmhdbprd02'],
  ['ddcb7992-2ffb-5c10-8b39-80685f6eaaba', 'vmdrbdqas01'],
  ['422686d6-b2d1-5092-93e8-a744854f5085', 'vmdrbdqas02'],
  ['25677e37-fd33-5005-896c-9275b1284534', 'vmnwqas01'],
  ['3711ea88-9ccc-5b07-8f9d-042be449d72b', 'vmnwqas02'],
  ['098fc159-3ed6-58e7-91be-38fda8a833ea', 'vmnwqas03'],
  ['81e9b629-c1e7-538f-bff1-47d3a6580522', 'vmnwqas04'],
  ['99cf8a3a-48d6-57a4-b302-6e4482227ab6', 'vmhdbqas01'],
  ['e0c182db-32ff-55c6-a9eb-2b82dd21bc8b', 'vmhdbqas02'],
  ['240f96b1-8d26-53b7-9e99-ffb0f2e735bf', 'vmdrbddev01'],
  ['21de186a-e38f-5804-b643-7f4ef22fecfd', 'vmdrbddev02'],
  ['7269ee51-5007-5849-aaa7-7c4a98b0c9ce', 'vmnwdev01'],
  ['fb2c6b8a-9915-5969-a6b7-8b5a42de1971', 'vmnwdev02'],
  ['9a3ec76a-dd4f-5013-9cf0-5eb4cf89898f', 'vmnwdev03'],
  ['1b0e9297-97dd-55d6-9874-8efde4d84c90', 'vmnwdev04'],
  ['13e8c25c-3180-5a9a-95c8-51ec38e50cfc', 'vmhdbdev01'],
  ['0a055c90-4cb6-54ce-ac9c-ae3fedaf40d4', 'vmhdbdev02'],
  ['69f4dcbb-efa2-5a16-8bc8-01df7dbb7384', 'vmiscsi01'],
  ['f0c808b3-d869-5192-a944-20f66a6a8449', 'vmiscsi01'],
  ['9a26b6d0-6e72-597c-9fe5-152a6875f214', 'vmiscsi01'],
];

export const allHostNames = () =>
  availableHosts.map(([_, hostName]) => hostName);
export const agents = () => availableHosts.map(([agentId, _]) => agentId);
