using AutoMapper;
using DFDSServiceAPI.Dtos;
using Microsoft.AspNetCore.Mvc;
using Service;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace DFDSServiceAPI.Controllers
{
    [Produces("application/json")]
    [Route("api/[controller]")]
    [ApiController]
    public class ServiceProxyController : Controller
    {
        private readonly IServiceProxyService _proxyService;
        private readonly IMapper _mapper;

        public ServiceProxyController(IServiceProxyService proxyService, IMapper mapper)
        {
            _proxyService = proxyService;
            _mapper = mapper;
        }

        [HttpGet]
        [ProducesResponseType(200)]
        [ProducesResponseType(404)]
        public async Task<IActionResult> GetAll()
        {
            var results = _mapper.Map<List<ServiceProxyResultDto>>(await _proxyService.GetResults());
            return Ok(results);
        }
    }
}
