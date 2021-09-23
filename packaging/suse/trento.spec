#
# spec file for package trento
#
# Copyright (c) 2021 SUSE LLC
#
# All modifications and additions to the file contributed by third parties
# remain the property of their copyright owners, unless otherwise agreed
# upon. The license for this file, and modifications and additions to the
# file, is the same license as for the pristine package itself (unless the
# license for the pristine package is not an Open Source License, in which
# case the license is the MIT License). An "Open Source License" is a
# license that conforms to the Open Source Definition (Version 1.9)
# published by the Open Source Initiative   .

# Please submit bugfixes or comments via https://bugs.opensuse.org/
#


Name:           trento
# Version will be processed via set_version source service
Version:        0
Release:        0
License:        Apache-2.0
Summary:        An open cloud-native web console improving on the life of SAP Applications administrators.
Group:          System/Monitoring
URL:            https://github.com/trento-project/trento
Source:         %{name}-%{version}.tar.gz
Source1:        vendor.tar.gz
Source2:        node_modules.spec.inc
Source3:	package.json
%include %_sourcedir/node_modules.spec.inc
ExclusiveArch:  aarch64 x86_64 ppc64le s390x
BuildRoot:      %{_tmppath}/%{name}-%{version}-build
BuildRequires:  golang-packaging
BuildRequires:  golang(API) = 1.16
BuildRequires:  npm
BuildRequires:  local-npm-registry
Provides:       trento = %{version}-%{release}

%{go_nostrip}

%description
An open cloud-native web console improving on the life of SAP Applications administrators.

Trento is a city on the Adige River in Trentino-Alto Adige/Suedtirol in Italy. [...] It is one of the nation's wealthiest and most prosperous cities, [...] often ranking highly among Italian cities for quality of life, standard of living, and business and job opportunities. (source)

This project is a reboot of the "SUSE Console for SAP Applications", also known as the Blue Horizon for SAP prototype, which is focused on automated infrastructure deployment and provisioning for SAP Applications.

As opposed to that first iteration, this new one will focus more on operations of existing clusters, rather than deploying new one.

%prep
%setup -q            # unpack project sources
%setup -q -T -D -a 1 # unpack go dependencies in vendor.tar.gz, which was prepared by the source services

cp %SOURCE3 .
local-npm-registry %{_sourcedir} install --with=dev

%define shortname trento

%build

mv node_modules web/frontend/
VERSION=%{version} make build

%install

# Install the binary.
install -D -m 0755 %{shortname} "%{buildroot}%{_bindir}/%{shortname}"

# Install the systemd unit
install -D -m 0644 trento-agent.service %{buildroot}%{_unitdir}/trento-agent.service

%pre
%service_add_pre trento-agent.service

%post
%service_add_post trento-agent.service

%preun
%service_del_preun trento-agent.service

%postun
%service_del_postun trento-agent.service

%files
%defattr(-,root,root)
%doc *.md
%doc docs/*.md
%license LICENSE
%{_bindir}/%{shortname}
%{_unitdir}/trento-agent.service

%changelog
