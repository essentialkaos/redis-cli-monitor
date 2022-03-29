################################################################################

# rpmbuilder:relative-pack true

################################################################################

%define  debug_package %{nil}

################################################################################

Summary:         Tiny Redis client for renamed MONITOR commands
Name:            redis-cli-monitor
Version:         2.2.1
Release:         0%{?dist}
Group:           Applications/System
License:         Apache License, Version 2.0
URL:             https://kaos.sh/redis-cli-monitor

Source0:         https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:       %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:   golang >= 1.17

Provides:        %{name} = %{version}-%{release}

################################################################################

%description
Tiny Redis client for renamed MONITOR commands.

################################################################################

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

%post
if [[ -d %{_sysconfdir}/bash_completion.d ]] ; then
  %{name} --completion=bash 1> %{_sysconfdir}/bash_completion.d/%{name} 2>/dev/null
fi

if [[ -d %{_datarootdir}/fish/vendor_completions.d ]] ; then
  %{name} --completion=fish 1> %{_datarootdir}/fish/vendor_completions.d/%{name}.fish 2>/dev/null
fi

if [[ -d %{_datadir}/zsh/site-functions ]] ; then
  %{name} --completion=zsh 1> %{_datadir}/zsh/site-functions/_%{name} 2>/dev/null
fi

%postun
if [[ $1 == 0 ]] ; then
  if [[ -f %{_sysconfdir}/bash_completion.d/%{name} ]] ; then
    rm -f %{_sysconfdir}/bash_completion.d/%{name} &>/dev/null || :
  fi

  if [[ -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish ]] ; then
    rm -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish &>/dev/null || :
  fi

  if [[ -f %{_datadir}/zsh/site-functions/_%{name} ]] ; then
    rm -f %{_datadir}/zsh/site-functions/_%{name} &>/dev/null || :
  fi
fi

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_bindir}/%{name}

################################################################################

%changelog
* Wed Mar 30 2022 Anton Novojilov <andy@essentialkaos.com> - 2.2.1-0
- Removed pkg.re usage
- Added module info
- Added Dependabot configuration

* Tue Sep 22 2020 Anton Novojilov <andy@essentialkaos.com> - 2.2.0-0
- Added option for filtering data by DB number
- ek package updated to the latest stable version

* Thu Oct 17 2019 Anton Novojilov <andy@essentialkaos.com> - 2.1.1-0
- ek package updated to the latest stable version

* Fri Jun 14 2019 Anton Novojilov <andy@essentialkaos.com> - 2.1.0-0
- ek package updated to the latest stable version
- Added completion generation for bash, zsh and fish

* Sat Oct 20 2018 Anton Novojilov <andy@essentialkaos.com> - 2.0.2-0
- Show usage info if '-h' passed without any value

* Thu Jul 06 2017 Anton Novojilov <andy@essentialkaos.com> - 2.0.1-0
- Fixed bug with handling Redis errors

* Sun Jul 02 2017 Anton Novojilov <andy@essentialkaos.com> - 2.0.0-0
- Added colors and timestamp formatting
- Added option for enabling raw output
- Code refactoring
- Fixed bug in usage examples

* Wed Jun 07 2017 Anton Novojilov <andy@essentialkaos.com> - 1.4.0-0
- Minor improvements

* Fri May 26 2017 Anton Novojilov <andy@essentialkaos.com> - 1.3.0-0
- ek package updated to v9

* Sun Apr 16 2017 Anton Novojilov <andy@essentialkaos.com> - 1.2.0-0
- ek package updated to v8

* Fri Mar 10 2017 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-0
- ek package updated to v7

* Tue Oct 11 2016 Anton Novojilov <andy@essentialkaos.com> - 1.0.5-0
- ek package updated to v5

* Fri Sep 16 2016 Anton Novojilov <andy@essentialkaos.com> - 1.0.4-0
- ek package updated to v3

* Fri Jun 05 2015 Anton Novojilov <andy@essentialkaos.com> - 1.0.1-0
- Fixed bug with arguments parsing

* Tue Apr 21 2015 Anton Novojilov <andy@essentialkaos.com> - 1.0-1
- Fixed description

* Wed Mar 11 2015 Anton Novojilov <andy@essentialkaos.com> - 1.0-0
- Initial build
