using AutoMapper;
using DFDSServiceAPI.Dtos;
using Service.Classes;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace DFDSServiceAPI
{
    public class AutomapperProfile : Profile
    {
        public AutomapperProfile()
        {
            CreateMap<ServiceProxyResult, ServiceProxyResultDto>();
        }

    }
}
