###############################################################################

# rpmbuilder:relative-pack true

###############################################################################

%define  debug_package %{nil}

###############################################################################

Summary:         Tiny redis client for renamed MONITOR commands
Name:            redis-cli-monitor
Version:         1.1.0
Release:         0%{?dist}
Group:           Applications/System
License:         EKOL
URL:             http://essentialkaos.com

Source0:         https://source.kaos.io/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:       %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:   golang >= 1.7

Provides:        %{name} = %{version}-%{release}

###############################################################################

%description
Tiny redis client for renamed MONITOR commands.

###############################################################################

%prep
%setup -q

%build
export GOPATH=$(pwd) 
go build src/github.com/essentialkaos/%{name}/%{name}.go

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -pm 755 %{name} %{buildroot}%{_bindir}/

%clean
rm -rf %{buildroot}

###############################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE.EN LICENSE.RU
%{_bindir}/%{name}

###############################################################################

%changelog
* Fri Mar 10 2017 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-0
- EK package updated to v7

* Tue Oct 11 2016 Anton Novojilov <andy@essentialkaos.com> - 1.0.5-0
- EK package updated to v5

* Fri Sep 16 2016 Anton Novojilov <andy@essentialkaos.com> - 1.0.4-0
- EK package updated to v3

* Fri Jun 05 2015 Anton Novojilov <andy@essentialkaos.com> - 1.0.1-0
- Fixed bug with arguments parsing

* Tue Apr 21 2015 Anton Novojilov <andy@essentialkaos.com> - 1.0-1
- Fixed description

* Wed Mar 11 2015 Anton Novojilov <andy@essentialkaos.com> - 1.0-0
- Initial build
